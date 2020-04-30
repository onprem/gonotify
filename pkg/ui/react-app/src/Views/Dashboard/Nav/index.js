import React from 'react';
import { Link } from 'react-router-dom';

import { useAuth } from 'Context/auth';

import { ReactComponent as SettingIcon } from 'Assets/icons/settings.svg';

import styles from './nav.module.css';

const Nav = ({ user }) => {
  const { logMeOut } = useAuth();
  const name = user?.name || 'Irfan';

  return (
    <div className={styles.nav}>
      <Link to="/dashboard/account" className={styles.icoLink}>
        <SettingIcon className={styles.ico} />
      </Link>
      <div className={styles.accDiv}>
        <img src={`https://api.adorable.io/avatars/60/${name.replace(/\s/g,'')}.png`} alt="avatar" className={styles.avatar} />
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
