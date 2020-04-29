import React from 'react';
import { Link } from 'react-router-dom';

import { useAuth } from 'Context/auth';

import user from 'Assets/icons/acc.svg';
import { ReactComponent as SettingIcon } from 'Assets/icons/settings.svg';

import styles from './nav.module.css';

const Nav = () => {
  const { logMeOut } = useAuth();
  return (
    <div className={styles.nav}>
      <Link to="/dashboard/account" className={styles.icoLink}>
        <SettingIcon className={styles.ico} />
      </Link>
      <div className={styles.accDiv}>
        <img src={user} alt="avatar" className={styles.avatar} />
        <div className={styles.menu}>
          <a className={styles.link} href="#logout" onClick={logMeOut}>
            LogOut
          </a>
        </div>
      </div>
    </div>
  );
};

export default Nav;
