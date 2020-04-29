import React from 'react';
import { Link, useParams } from 'react-router-dom';
import { Text } from 'Components/Form';
import Button from 'Components/Button';

import { ReactComponent as ArrowIcon } from 'Assets/icons/arrow.svg';
import styles from '../logreg.module.css';

const Verify = () => {
  const { phone } = useParams();

  return (
    <>
      <h1 className={styles.heading}>Lets join !!</h1>
      <p className={styles.para}>Verify your phone number to continue</p>
      <form className={styles.form}>
        <Text name="phone" label="Your Phone" value={phone} readOnly={true} />
        <Text type="text" name="code" label="Verification Code" placeholder="123456" />
        <Button className={styles.btn} type="submit">
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
