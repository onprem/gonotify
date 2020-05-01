import React from 'react';
import { NavLink } from 'react-router-dom';

import Button from 'Components/Button';

import { ReactComponent as MonoLogo } from 'Assets/logo/mono.svg';
import { ReactComponent as AccIcon } from 'Assets/icons/acc.svg';
import { ReactComponent as GroupIcon } from 'Assets/icons/group.svg';
import { ReactComponent as NumIcon } from 'Assets/icons/member.svg';
import { ReactComponent as SendIcon } from 'Assets/icons/send.svg';
import { ReactComponent as HelpIcon } from 'Assets/icons/help.svg';

import styles from './sidebar.module.css';

const Sidebar = () => {
  return (
    <section className={styles.sidebar}>
      <div className={styles.top}>
        <MonoLogo className={styles.logo} />
        <div className={styles.nav}>
          <NavLink activeClassName={styles.active} className={styles.link} to="/dashboard/account">
            <AccIcon className={styles.icon} /> <span className={styles.linkTitle}>Account</span>
          </NavLink>
          <NavLink activeClassName={styles.active} className={styles.link} to="/dashboard/groups">
            <GroupIcon className={styles.icon} /> <span className={styles.linkTitle}>Groups</span>
          </NavLink>
          <NavLink activeClassName={styles.active} className={styles.link} to="/dashboard/numbers">
            <NumIcon className={styles.icon} /> <span className={styles.linkTitle}>Numbers</span>
          </NavLink>
          <NavLink activeClassName={styles.active} className={styles.link} to="/dashboard/notifs">
            <SendIcon className={styles.icon} /> <span className={styles.linkTitle}>Notifications</span>
          </NavLink>
        </div>
      </div>
      <Button className={styles.help}>
        <HelpIcon className={styles.helpicon} /> <span className={styles.linkTitle}>Help</span>
      </Button>
    </section>
  );
};

export default Sidebar;
