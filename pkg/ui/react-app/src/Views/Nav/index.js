import React from 'react';
import { Link } from 'react-router-dom';
import Button from 'Components/Button';

import { ReactComponent as Logo } from 'Assets/logo/logo.svg';
import styles from './nav.module.css';

const Nav = () => {
  return (
    <nav className={styles.nav}>
      <Link to="/">
        <Logo className={styles.logo} />
      </Link>
      <div className={styles.links}>
        <Link className={styles.link} to="/about">
          About
        </Link>
        <Link className={styles.link} to="/contact">
          Contact
        </Link>
        <Link className={styles.link} to="/help">
          Help
        </Link>
      </div>
      <Link to="/login">
        <Button style={{ backgroundColor: 'var(--colorBrand-light)' }}>Sign in</Button>
      </Link>
    </nav>
  );
};

export default Nav;
