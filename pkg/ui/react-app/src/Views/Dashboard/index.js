import React from 'react';
import { Switch, Route } from 'react-router-dom';
import useSWR from 'swr';

import Sidebar from './Sidebar';
import Nav from './Nav';
import Account from 'Views/Account';
import Groups from 'Views/Groups';
import Numbers from 'Views/Numbers';
import GroupDetails from 'Views/GroupDetails';

import styles from './dashboard.module.css';

const Dashboard = () => {
  const { data: gdata, error: gerror } = useSWR('/api/v1/groups');
  const { data: udata, error: uerror } = useSWR('/api/v1/user');

  return (
    <div className={styles.dash}>
      <Sidebar />
      <section className={styles.main}>
        <Nav user={udata} />
        <header className={styles.header}>
          <h1 className={styles.heading}>Dashboard</h1>
        </header>
        {(gerror || uerror) ? (
          <h1>Some Error Occured</h1>
        ) : (gdata && udata) ? (
          <Switch>
            <Route exact path="/dashboard/account">
              <Account user={udata} />
            </Route>
            <Route exact path="/dashboard/groups">
              <Groups groups={gdata.groups} />
            </Route>
            <Route exact path="/dashboard/groups/:name">
              <GroupDetails groups={gdata.groups} />
            </Route>
            <Route exact path="/dashboard/numbers">
              <Numbers groups={gdata.groups} />
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
