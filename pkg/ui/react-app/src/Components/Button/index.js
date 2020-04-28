import React from 'react';
import classNames from 'classnames';
import styles from './button.module.css';

const Button = ({ size, style, className, children }) => {
  return <button style={style} className={classNames(styles.btn, styles[size], className)}>{children}</button>;
};

Button.defaultProps = {
  size: 'regular',
  style: undefined,
};

export default Button;
