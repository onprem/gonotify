import React, { useState, useEffect } from 'react';
import useSWR from 'swr';
import { useForm } from 'react-hook-form';

import Button from 'Components/Button';
import { Text } from 'Components/Form';
import Modal from 'Components/Modal';
import { ReactComponent as SendIcon } from 'Assets/icons/send.svg';
import useSWRPost from 'Hooks/useSWRPost';
import toast from 'Utils/toast';

import styles from './groups.module.css';

const SendMsg = ({ name, setIsOpen }) => {
  const { handleSubmit, register, errors } = useForm();
  const [sendMsg, { isValidating }] = useSWRPost('/api/v1/send/whatsapp', {
    onSuccess: (data) => {
      if (data.error) toast.error(data.error);
      else {
        toast.success(data.message);
        setIsOpen(false);
      }
    },
    onError: toast.error,
  });

  useEffect(() => {
    const fields = ['phone', 'password'];
    fields.forEach((f) => {
      if (errors[f]?.message) toast.error(errors[f].message);
    });
  });

  return (
    <form className={styles.form} onSubmit={handleSubmit(sendMsg)}>
      <input
        type="hidden"
        name="name"
        value={name}
        ref={register({
          required: 'Group name is required',
        })}
      />
      <Text
        name="body"
        label="Message"
        placeholder="Hi there!"
        inpRef={register({
          required: 'Message body is required',
        })}
      />
      <Button type="submit" className={styles.formBtn} disabled={isValidating}>
        SEND <SendIcon />
      </Button>
    </form>
  );
};

const Card = ({ name, id, groups }) => {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <>
      <div className={styles.card}>
        <div className={styles.up}>
          <h3 className={styles.cardName}>{name}</h3>
          <p className={styles.cardDetail}>ID: {id}</p>
          <p className={styles.cardDetail}>Members: {groups}</p>
        </div>
        <Button className={styles.btn} onClick={() => setIsOpen(true)}>
          SEND <SendIcon />
        </Button>
      </div>
      <Modal isOpen={isOpen} setIsOpen={setIsOpen}>
        <SendMsg name={name} setIsOpen={setIsOpen} />
      </Modal>
    </>
  );
};

const Groups = () => {
  const { data, error } = useSWR('/api/v1/groups');

  if (error) return <h1>Some error occured.</h1>;
  if (!data) return <h1>Loading...</h1>;

  return (
    <div className={styles.groups}>
      <h2 className={styles.heading}>All Groups</h2>
      <div className={styles.cards}>
        {data.groups.map((g) => (
          <Card key={g.id} name={g.name} id={g.id} groups={g.whatsappNodes.length} />
        ))}
      </div>
    </div>
  );
};

export default Groups;
