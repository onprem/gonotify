import React from 'react';

import styles from './modal.module.css';

const Modal = ({ isOpen, setIsOpen, children }) => {
  if (!isOpen) return null;
  return (
    <div className={styles.container} onClick={() => setIsOpen(false)}>
      <div className={styles.inner} onClick={(e) => e.stopPropagation()}>
        {children}
      </div>
    </div>
  );
};

export default Modal;
