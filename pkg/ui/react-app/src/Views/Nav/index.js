import React from 'react';
import Button from 'Components/Button';

import { ReactComponent as Logo } from 'Assets/logo/logo.svg';
import styles from './nav.module.css';

const Nav = () => {
  return (
    <nav className={styles.nav}>
      <Logo className={styles.logo} />
      <div className={styles.links}>
        <a className={styles.link} href="#about">About</a>
        <a className={styles.link} href="#contact">Contact</a>
        <a className={styles.link} href="#help">Help</a>
      </div>
      <Button style={{ backgroundColor: 'var(--colorBrand-light)' }}>Sign in</Button>
    </nav>
  );
};

export default Nav;
