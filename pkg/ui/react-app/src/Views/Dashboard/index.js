import React from 'react';
import Sidebar from './Sidebar';
import Nav from './Nav';

import styles from './dashboard.module.css';

const Dashboard = () => {
  return (
    <div className={styles.dash}>
      <Sidebar />
      <section className={styles.main}>
        <Nav />
        <header className={styles.header}>
          <h1 className={styles.heading}>Dashboard</h1>
        </header>
      </section>
    </div>
  );
};

export default Dashboard;
