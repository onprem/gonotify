import React from 'react';
import Label from '../Label';
import styles from './text.module.css';

const Text = ({ className, type, label, ...rest }) => {
  return (
    <Label text={label}>
      <input {...rest} type={type} className={styles.input} />
    </Label>
  );
};

Text.defaultProps = {
  type: 'text',
};

export default Text;
