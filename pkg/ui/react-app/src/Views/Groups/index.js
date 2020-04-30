import React, { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { Link } from 'react-router-dom';
import { mutate } from 'swr';

import Button from 'Components/Button';
import { Text } from 'Components/Form';
import Modal from 'Components/Modal';
import { ReactComponent as SendIcon } from 'Assets/icons/send.svg';
import { ReactComponent as PlusIcon } from 'Assets/icons/plus.svg';
import { ReactComponent as TrashIcon } from 'Assets/icons/trash.svg';
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

const RemoveGroup = ({ groupID, setIsOpen }) => {
  const { handleSubmit, register, errors } = useForm();
  const [removeGroup, { isValidating }] = useSWRPost('/api/v1/groups/remove', {
    onSuccess: (data) => {
      if (data.error) toast.error(data.error);
      else {
        toast.success(data.message);
        mutate('/api/v1/numbers');
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
  });

  const onSubmit = (v) => {
    removeGroup({
      id: Number(v.id),
    });
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit(onSubmit)}>
      <input
        type="hidden"
        name="id"
        value={groupID}
        ref={register({
          required: 'GroupID is required',
        })}
      />
      <h3>This action is permanent.</h3>
      <h3>Are you sure?</h3>
      <Button type="submit" className={styles.delBtn} disabled={isValidating}>
        Remove <TrashIcon />
      </Button>
    </form>
  );
};

const Card = ({ name, id, groups }) => {
  const [isOpen, setIsOpen] = useState(false);
  const [isDelOpen, setDelOpen] = useState(false);

  return (
    <>
      <div className={styles.card}>
        <div className={styles.up}>
          <div className={styles.ph}>
            <Link to={`/dashboard/groups/${name}`}>
              <h3 className={styles.cardName}>{name}</h3>
            </Link>
            <Button className={styles.trashBtn} onClick={() => setDelOpen(true)}>
              <TrashIcon />
            </Button>
          </div>
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
      <Modal isOpen={isDelOpen} setIsOpen={setDelOpen}>
        <RemoveGroup groupID={id} setIsOpen={setDelOpen} />
      </Modal>
    </>
  );
};

const NewGroup = ({ setIsOpen }) => {
  const { handleSubmit, register, errors } = useForm();
  const [addGroup, { isValidating }] = useSWRPost('/api/v1/groups/add', {
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
    const fields = ['name'];
    fields.forEach((f) => {
      if (errors[f]?.message) toast.error(errors[f].message);
    });
  });

  return (
    <form className={styles.form} onSubmit={handleSubmit(addGroup)}>
      <Text
        name="name"
        label="Name"
        placeholder="ex. workplace"
        inpRef={register({
          required: 'Name is required',
          pattern: {
            value: /^[a-zA-Z]+$/,
            message: 'Only alphabetic characters are allowed',
          },
        })}
      />
      <Button type="submit" className={styles.formBtn} disabled={isValidating}>
        ADD <PlusIcon />
      </Button>
    </form>
  );
};

const Groups = ({ groups }) => {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <>
      <div className={styles.groups}>
        <h2 className={styles.heading}>All Groups</h2>
        <div className={styles.cards}>
          {groups.map((g) => (
            <Card key={g.id} name={g.name} id={g.id} groups={g.whatsappNodes?.length || 0} />
          ))}
          <div className={styles.card}>
            <Button className={styles.addBtn} onClick={() => setIsOpen(true)}>
              <PlusIcon /> Add Group
            </Button>
          </div>
        </div>
      </div>
      <Modal isOpen={isOpen} setIsOpen={setIsOpen}>
        <NewGroup setIsOpen={setIsOpen} />
      </Modal>
    </>
  );
};

export default Groups;
