import React, { useState, useEffect } from 'react';
import useSWR, { mutate } from 'swr';
import { useForm } from 'react-hook-form';

import Button from 'Components/Button';
import Modal from 'Components/Modal';
import { Select } from 'Components/Form';
import { ReactComponent as PlusIcon } from 'Assets/icons/plus.svg';
import useSWRPost from 'Hooks/useSWRPost';
import toast from 'Utils/toast';

import styles from './numbers.module.css';

const AddToGroup = ({ numberID, groups, setIsOpen }) => {
  const { handleSubmit, register, errors } = useForm();
  const [sendMsg, { isValidating }] = useSWRPost('/api/v1/whatsapps/group/add', {
    onSuccess: (data) => {
      if (data.error) toast.error(data.error);
      else {
        toast.success(data.message);
        mutate('/api/v1/numbers');
        mutate('/api/v1/groups')
        setIsOpen(false);
      }
    },
    onError: toast.error,
  });

  useEffect(() => {
    const fields = ['groupID', 'numberID'];
    fields.forEach((f) => {
      if (errors[f]?.message) toast.error(errors[f].message);
    });
  });

  const onSubmit = (v) => {
    sendMsg({
      groupID: Number(v.groupID),
      numberID: Number(v.numberID),
    });
  };

  const filter = groups.filter((g) => !g.whatsappNodes.some((n) => n.numberID === numberID));

  if (filter.length === 0) return <h3 style={{ margin: '1em 2em'}}>This number is already in all groups.</h3>;

  return (
    <form className={styles.form} onSubmit={handleSubmit(onSubmit)}>
      <input
        type="hidden"
        name="numberID"
        value={numberID}
        ref={register({
          required: 'NumberID is required',
        })}
      />
      <Select
        name="groupID"
        label="Group"
        inpRef={register({
          required: 'GroupID is required',
        })}
        options={filter.map((g) => ({ value: g.id, text: g.name }))}
      />
      <Button type="submit" className={styles.formBtn} disabled={isValidating}>
        ADD <PlusIcon />
      </Button>
    </form>
  );
};

const Card = ({ phone, id, groups, allGroups }) => {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <>
      <div className={styles.card}>
        <div className={styles.content}>
          <h3 className={styles.cardName}>{phone}</h3>
          <p className={styles.cardDetail}>ID: {id}</p>
          <p className={styles.cardDetail}>Groups: {groups}</p>
        </div>
        <Button className={styles.btn} onClick={() => setIsOpen(true)}>
          Add to Group <PlusIcon />
        </Button>
      </div>
      <Modal isOpen={isOpen} setIsOpen={setIsOpen}>
        <AddToGroup numberID={id} groups={allGroups} setIsOpen={setIsOpen} />
      </Modal>
    </>
  );
};

const Numbers = ({ groups }) => {
  const { data, error } = useSWR('/api/v1/numbers');

  if (error) return <h1>Some error occured.</h1>;
  if (!data) return <h1>Loading...</h1>;
  return (
    <div className={styles.numbers}>
      <h2 className={styles.heading}>All Numbers</h2>
      <div className={styles.cards}>
        {data.numbers.map((n) => (
          <Card key={n.id} id={n.id} phone={n.phone} groups={n.groups} allGroups={groups} />
        ))}
      </div>
    </div>
  );
};

export default Numbers;
