import React from 'react';
import { Text } from 'Components/Form'

import styles from './acc.module.css';

const Account = ({ user }) => {
  return (
    <div className={styles.container}>
      <div className={styles.acc}>
      <h2 className={styles.heading}>Account Section</h2>
      <div>
        <Text name="name" label="Name" value={user.name} readOnly />
        <Text name="phone" label="Phone number" value={user.phone} readOnly />
      </div>
      </div>
    </div>
  );
};

export default Account;
