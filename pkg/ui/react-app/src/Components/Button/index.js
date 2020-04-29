import React from 'react';
import classNames from 'classnames';
import styles from './button.module.css';

const Button = ({ size, style, type, className, disabled, onClick, children }) => {
  return (
    <button
      type={type}
      style={style}
      className={classNames(styles.btn, styles[size], className)}
      disabled={disabled}
      onClick={onClick}
    >
      {children}
    </button>
  );
};

Button.defaultProps = {
  size: 'regular',
  style: undefined,
  type: 'button',
  disabled: false,
};

export default Button;
