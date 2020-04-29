import React from 'react';
import { Switch, Route } from 'react-router-dom';

import Sidebar from './Sidebar';
import Nav from './Nav';
import Groups from 'Views/Groups';

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
        <Switch>
          <Route exact path="/dashboard/groups">
            <Groups />
          </Route>
        </Switch>
      </section>
    </div>
  );
};

export default Dashboard;
