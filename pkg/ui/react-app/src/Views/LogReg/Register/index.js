import React, { useEffect } from 'react';
import { Link, useHistory } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import useSWRPost from 'Hooks/useSWRPost';

import { Text } from 'Components/Form';
import Button from 'Components/Button';

import toast from 'Utils/toast';

import { ReactComponent as ArrowIcon } from 'Assets/icons/arrow.svg';
import styles from './register.module.css';

const Register = () => {
  const { handleSubmit, register, errors, getValues } = useForm();
  const history = useHistory();

  const [runRegister, { isValidating }] = useSWRPost('/api/v1/register', {
    onSuccess: (data) => {
      if (data.error) toast.error(data.error);
      else {
        toast.success(data.message);
        history.push(`/verify/${getValues('phone')}`);
      }
    },
    onError: toast.error,
  });

  useEffect(() => {
    const fields = ['name', 'phone', 'password'];
    fields.forEach((f) => {
      if (errors[f]?.message) toast.error(errors[f].message);
    });
  });

  return (
    <>
      <h1 className={styles.heading}>Lets join !!</h1>
      <p className={styles.para}>Enter name, phone number and password to continue</p>
      <form className={styles.form} onSubmit={handleSubmit(runRegister)}>
        <Text
          name="name"
          label="Full Name"
          placeholder="Irfan Khan"
          inpRef={register({
            required: 'Name is required',
          })}
        />
        <Text
          name="phone"
          label="Your Phone"
          placeholder="+919912312345"
          inpRef={register({
            required: 'Phone is required',
            pattern: {
              value: /^\+?\d+$/,
              message: 'Invalid phone number',
            },
          })}
        />
        <Text
          type="password"
          name="password"
          label="Password"
          placeholder="password"
          inpRef={register({
            required: 'Password is required',
            minLength: {
              value: 4,
              message: 'Minimum password length is 4',
            },
          })}
        />
        <Button className={styles.btn} type="submit" disabled={isValidating}>
          Register <ArrowIcon />
        </Button>
      </form>
      <hr className={styles.hr} />
      <p className={styles.para}>
        Already have an account? <Link to="/login">Sign In</Link>
      </p>
    </>
  );
};

export default Register;
