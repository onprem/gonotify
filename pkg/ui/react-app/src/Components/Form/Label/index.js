import React from 'react';

import styles from './label.module.css';

const Label = ({ text, children}) => {
  return (
    <label className={styles.label}>
      <span className={styles.text}>{text}</span>
      {children}
    </label>
  );
};

export default Label;
