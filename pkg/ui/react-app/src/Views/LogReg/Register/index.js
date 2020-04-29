import React from 'react';
import { Link } from 'react-router-dom';
import { Text } from 'Components/Form';
import Button from 'Components/Button';

import { ReactComponent as ArrowIcon } from 'Assets/icons/arrow.svg';
import styles from './register.module.css';

const Register = () => {
  return (
    <>
      <h1 className={styles.heading}>Lets join !!</h1>
      <p className={styles.para}>Enter name, phone number and password to continue</p>
      <from className={styles.form}>
        <Text name="name" label="Full Name" placeholder="Irfan Khan" />
        <Text name="phone" label="Your Phone" placeholder="+919912312345" />
        <Text type="password" name="password" label="Password" placeholder="password" />
        <Button className={styles.btn} type="submit">
          Register <ArrowIcon />
        </Button>
      </from>
      <hr className={styles.hr} />
      <p className={styles.para}>
        Already have an account? <Link to="/login">Sign In</Link>
      </p>
    </>
  );
};

export default Register;
