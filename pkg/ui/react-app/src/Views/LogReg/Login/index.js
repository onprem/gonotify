import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import useSWR from 'swr';
import toast from 'Utils/toast';
import { useAuth } from 'Context/auth';

import { Text } from 'Components/Form';
import Button from 'Components/Button';

import { ReactComponent as ArrowIcon } from 'Assets/icons/arrow.svg';
import styles from './login.module.css';

const Login = () => {
  const [values, setValues] = useState();
  const { handleSubmit, register, errors } = useForm();
  const { setToken } = useAuth();

  const { isValidating } = useSWR(values ? ['/api/v1/login', 'POST', values] : null, {
    onSuccess: (data) => {
      if (data.error) toast.error(data.error);
      else {
        setToken(data.token);
        toast.success(data.message);
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

  const onSubmit = setValues;

  return (
    <>
      <h1 className={styles.heading}>Get in !!</h1>
      <p className={styles.para}>Enter your phone number and password to continue</p>
      <form className={styles.form} onSubmit={handleSubmit(onSubmit)}>
        <Text
          name="phone"
          label="Your Phone"
          placeholder="+919912312345"
          inpRef={register({
            required: 'Phone number is required',
            pattern: {
              value: /^\+?\d+$/,
              message: 'Invalid phone number',
            },
          })}
          errored={errors.phone}
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
          errored={errors.password}
        />
        <Button className={styles.btn} type="submit" disabled={isValidating}>
          Sign In <ArrowIcon />
        </Button>
      </form>
      <hr className={styles.hr} />
      <p className={styles.para}>
        Don't have an account yet? <Link to="/register">Register</Link>
      </p>
    </>
  );
};

export default Login;
