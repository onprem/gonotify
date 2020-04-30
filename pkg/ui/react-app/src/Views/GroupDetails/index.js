import React, {useEffect, useState} from 'react';
import { useParams } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { mutate } from 'swr'

import Button from 'Components/Button';
import Modal from 'Components/Modal';
import { ReactComponent as TrashIcon } from 'Assets/icons/trash.svg';
import useSWRPost from 'Hooks/useSWRPost';
import toast from 'Utils/toast';

import styles from './groupdetails.module.css';

const RemoveFromGroup = ({ id, setIsOpen }) => {
  const { handleSubmit, register, errors } = useForm();
  const [sendMsg, { isValidating }] = useSWRPost('/api/v1/whatsapps/group/remove', {
    onSuccess: (data) => {
      if (data.error) toast.error(data.error);
      else {
        toast.success(data.message);
        mutate('/api/v1/groups');
        setIsOpen(false);
      }
    },
    onError: toast.error,
  });

  useEffect(() => {
    const fields = ['id'];
    fields.forEach((f) => {
      if (errors[f]?.message) toast.error(errors[f].message);
    });
  }, [errors]);

  const onSubmit = (v) => {
    sendMsg({
      id: Number(v.id),
    });
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit(onSubmit)}>
      <input
        type="hidden"
        name="id"
        value={id}
        ref={register({
          required: 'ID is required',
        })}
      />
      <h3>Are you sure?</h3>
      <Button type="submit" className={styles.formBtn} disabled={isValidating}>
        Remove from Group <TrashIcon />
      </Button>
    </form>
  );
};

const Card = ({ node }) => {
  const [isOpen, setIsOpen] = useState(false);
  const { id, phone, numberID } = node;

  return (
    <>
    <div className={styles.card}>
      <div className={styles.content}>
        <h3 className={styles.cardName}>{phone}</h3>
        <p className={styles.cardDetail}>ID: {id}</p>
        <p className={styles.cardDetail}>NumberID: {numberID}</p>
      </div>
      <Button className={styles.btn} onClick={() => setIsOpen(true)}>
        Remove <TrashIcon />
      </Button>
    </div>
    <Modal isOpen={isOpen} setIsOpen={setIsOpen}><RemoveFromGroup id={id} setIsOpen={setIsOpen} /></Modal>
    </>
  );
};

const GroupDetails = ({ groups }) => {
  const { name } = useParams();

  const [group] = groups.filter((g) => g.name === name.toLowerCase());

  if (!group) return <h2>Invalid Group name - {name}</h2>;

  return (
    <div className={styles.main}>
      <h2 className={styles.heading}>{group.name}</h2>
      <h3 className={styles.sub}>ID:&nbsp; {group.id}</h3>
      <div className={styles.cards}>
        {group.whatsappNodes.map((n) => (
          <Card key={n.id} node={n} />
        ))}
      </div>
    </div>
  );
};

export default GroupDetails;
