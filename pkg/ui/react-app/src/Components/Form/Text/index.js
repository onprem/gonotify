import React from 'react';
import classNames from 'classnames';

import Label from '../Label';
import styles from './text.module.css';

const Text = ({ className, type, label, inpRef, errored, ...rest }) => {
  return (
    <Label text={label}>
      <input
        {...rest}
        type={type}
        className={classNames(styles.input, { [styles.error]: errored })}
        ref={inpRef}
      />
    </Label>
  );
};

Text.defaultProps = {
  type: 'text',
  errored: false,
};

export default Text;
