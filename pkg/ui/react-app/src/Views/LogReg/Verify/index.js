import React, { useEffect } from 'react';
import { Link, useParams, useHistory } from 'react-router-dom';
import { useForm } from 'react-hook-form';

import useSWRPost from 'Hooks/useSWRPost';
import toast from 'Utils/toast';

import { Text } from 'Components/Form';
import Button from 'Components/Button';

import { ReactComponent as ArrowIcon } from 'Assets/icons/arrow.svg';
import styles from '../logreg.module.css';

const Verify = () => {
  const { phone } = useParams();
  const history = useHistory();
  const { handleSubmit, register, errors } = useForm();

  const [runVerify, { isValidating }] = useSWRPost('/api/v1/verify', {
    onSuccess: (data) => {
      if (data.error) toast.error(data.error);
      else {
        toast.success(data.message);
        history.push(`/login`);
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
    <>
      <h1 className={styles.heading}>Lets join !!</h1>
      <p className={styles.para}>Verify your phone number to continue</p>
      <form className={styles.form} onSubmit={handleSubmit(runVerify)}>
        <Text
          name="phone"
          label="Your Phone"
          value={phone}
          readOnly={true}
          inpRef={register({
            required: 'Phone number is required',
            pattern: {
              value: /^\+?\d+$/,
              message: 'Invalid phone number',
            },
          })}
        />
        <Text
          name="code"
          label="Verification Code"
          placeholder="123456"
          inpRef={register({
            required: 'Verification code is required',
            pattern: {
              value: /^\d{6}$/,
              message: 'Invalid verification code',
            },
          })}
        />
        <Button className={styles.btn} type="submit" disabled={isValidating}>
          Verify <ArrowIcon />
        </Button>
      </form>
      <hr className={styles.hr} />
      <p className={styles.para}>
        Already have an account? <Link to="/login">Sign In</Link>
      </p>
    </>
  );
};

export default Verify;
