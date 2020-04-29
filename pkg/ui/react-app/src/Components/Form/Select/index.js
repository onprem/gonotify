import React from 'react';
import classNames from 'classnames';

import Label from '../Label';
import styles from './select.module.css';

const Select = ({ name, className, label, inpRef, errored, options, ...rest }) => {
  return (
    <Label text={label}>
      <select
        {...rest}
        name={name}
        className={classNames(styles.input, { [styles.error]: errored })}
        ref={inpRef}
      >
        {options.map((o) => (
          <option key={o.value} value={o.value} className={styles.option}>
            {o.text}
          </option>
        ))}
      </select>
    </Label>
  );
};

Select.defaultProps = {
  errored: false,
};

export default Select;
