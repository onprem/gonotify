import React from 'react';
import { Switch, Route } from 'react-router-dom';
import useSWR from 'swr';

import Sidebar from './Sidebar';
import Nav from './Nav';
import Groups from 'Views/Groups';
import Numbers from 'Views/Numbers';
import GroupDetails from 'Views/GroupDetails';

import styles from './dashboard.module.css';

const Dashboard = () => {
  const { data, error } = useSWR('/api/v1/groups');

  return (
    <div className={styles.dash}>
      <Sidebar />
      <section className={styles.main}>
        <Nav />
        <header className={styles.header}>
          <h1 className={styles.heading}>Dashboard</h1>
        </header>
        {error ? (
          <h1>Some Error Occured</h1>
        ) : data ? (
          <Switch>
            <Route exact path="/dashboard/groups">
              <Groups groups={data.groups} />
            </Route>
            <Route exact path="/dashboard/groups/:name">
              <GroupDetails groups={data.groups} />
            </Route>
            <Route exact path="/dashboard/numbers">
              <Numbers groups={data.groups} />
            </Route>
          </Switch>
        ) : (
          <h1>Loading...</h1>
        )}
      </section>
    </div>
  );
};

export default Dashboard;
