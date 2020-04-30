import React, { useState, useEffect } from 'react';
import useSWR, { mutate } from 'swr';
import { useForm } from 'react-hook-form';

import Button from 'Components/Button';
import Modal from 'Components/Modal';
import { Select, Text } from 'Components/Form';
import { ReactComponent as PlusIcon } from 'Assets/icons/plus.svg';
import { ReactComponent as ArrowIcon } from 'Assets/icons/arrow.svg';
import { ReactComponent as TrashIcon } from 'Assets/icons/trash.svg';
import useSWRPost from 'Hooks/useSWRPost';
import toast from 'Utils/toast';

import styles from './numbers.module.css';

const AddToGroup = ({ numberID, groups, setIsOpen }) => {
  const { handleSubmit, register, errors } = useForm();
  const [addToGroup, { isValidating }] = useSWRPost('/api/v1/whatsapps/group/add', {
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
    const fields = ['groupID', 'numberID'];
    fields.forEach((f) => {
      if (errors[f]?.message) toast.error(errors[f].message);
    });
  });

  const onSubmit = (v) => {
    addToGroup({
      groupID: Number(v.groupID),
      numberID: Number(v.numberID),
    });
  };

  const filter = groups.filter((g) => !g.whatsappNodes.some((n) => n.numberID === numberID));

  if (filter.length === 0)
    return <h3 style={{ margin: '1em 2em' }}>This number is already in all groups.</h3>;

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

const RemoveNumber = ({ numberID, setIsOpen }) => {
  const { handleSubmit, register, errors } = useForm();
  const [removeNum, { isValidating }] = useSWRPost('/api/v1/numbers/remove', {
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
    const fields = ['numberID'];
    fields.forEach((f) => {
      if (errors[f]?.message) toast.error(errors[f].message);
    });
  });

  const onSubmit = (v) => {
    removeNum({
      numberID: Number(v.numberID),
    });
  };

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
      <h3>This action is permanent.</h3>
      <h3>Are you sure?</h3>
      <Button type="submit" className={styles.delBtn} disabled={isValidating}>
        Remove <TrashIcon />
      </Button>
    </form>
  );
};

const Card = ({ phone, id, groups, allGroups }) => {
  const [isOpen, setIsOpen] = useState(false);
  const [isRemOpen, setRemOpen] = useState(false);

  return (
    <>
      <div className={styles.card}>
        <div className={styles.content}>
          <div className={styles.ph}>
            <h3 className={styles.cardName}>{phone}</h3>
            <Button className={styles.trashBtn} onClick={() => setRemOpen(true)}>
              <TrashIcon />
            </Button>
          </div>
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
      <Modal isOpen={isRemOpen} setIsOpen={setRemOpen}>
        <RemoveNumber numberID={id} setIsOpen={setRemOpen} />
      </Modal>
    </>
  );
};

const VerifyNumber = ({ phone, setIsOpen }) => {
  const { handleSubmit, register, errors } = useForm();
  const [verifyNumber, { isValidating }] = useSWRPost('/api/v1/numbers/verify', {
    onSuccess: (data) => {
      if (data.error) toast.error(data.error);
      else {
        toast.success(data.message);
        mutate('/api/v1/numbers');
        setIsOpen(false);
      }
    },
    onError: toast.error,
  });

  useEffect(() => {
    const fields = ['phone', 'code'];
    fields.forEach((f) => {
      if (errors[f]?.message) toast.error(errors[f].message);
    });
  });

  return (
    <form className={styles.form} onSubmit={handleSubmit(verifyNumber)}>
      <input
        type="hidden"
        name="phone"
        value={phone}
        ref={register({
          required: 'Phone number is required',
        })}
      />
      <Text
        name="code"
        label="Verification Code"
        placeholder="123456"
        inpRef={register({
          required: 'Phone is required',
          pattern: {
            value: /^\d{6}$/,
            message: 'Invalid verification code',
          },
        })}
      />
      <Button type="submit" className={styles.formBtn} disabled={isValidating}>
        Verify <ArrowIcon />
      </Button>
    </form>
  );
};

const NewNumber = ({ setIsOpen }) => {
  const { handleSubmit, register, errors, getValues } = useForm();
  const [done, setDone] = useState(false);
  const [addNumber, { isValidating }] = useSWRPost('/api/v1/numbers/add', {
    onSuccess: (data) => {
      if (data.error) toast.error(data.error);
      else {
        toast.success(data.message);
        mutate('/api/v1/numbers');
        setDone(true);
      }
    },
    onError: toast.error,
  });

  useEffect(() => {
    const fields = ['phone'];
    fields.forEach((f) => {
      if (errors[f]?.message) toast.error(errors[f].message);
    });
  });

  if (done) return <VerifyNumber phone={getValues('phone')} setIsOpen={setIsOpen} />;

  return (
    <form className={styles.form} onSubmit={handleSubmit(addNumber)}>
      <Text
        name="phone"
        label="Phone"
        placeholder="+919912312345"
        inpRef={register({
          required: 'Phone is required',
          pattern: {
            value: /^\+?\d+$/,
            message: 'Invalid phone number',
          },
        })}
      />
      <Button type="submit" className={styles.formBtn} disabled={isValidating}>
        Add New Number <PlusIcon />
      </Button>
    </form>
  );
};

const Numbers = ({ groups }) => {
  const [isOpen, setIsOpen] = useState(false);
  const { data, error } = useSWR('/api/v1/numbers');

  if (error) return <h1>Some error occured.</h1>;
  if (!data) return <h1>Loading...</h1>;
  return (
    <>
      <div className={styles.numbers}>
        <h2 className={styles.heading}>All Numbers</h2>
        <div className={styles.cards}>
          {data.numbers.map((n) => (
            <Card key={n.id} id={n.id} phone={n.phone} groups={n.groups} allGroups={groups} />
          ))}
          <div className={styles.card}>
            <Button className={styles.addBtn} onClick={() => setIsOpen(true)}>
              <PlusIcon /> Add Number
            </Button>
          </div>
        </div>
      </div>
      <Modal isOpen={isOpen} setIsOpen={setIsOpen}>
        <NewNumber setIsOpen={setIsOpen} />
      </Modal>
    </>
  );
};

export default Numbers;
