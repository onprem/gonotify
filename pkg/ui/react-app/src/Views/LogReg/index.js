import React from 'react';
import { Switch, Route } from 'react-router-dom';

import Login from './Login';
import Register from './Register';

import { ReactComponent as BellIcon } from 'Assets/icons/notif.svg';

import styles from './logreg.module.css';

const LogReg = () => {
  return (
    <div className={styles.logreg}>
      <div className={styles.content}>
        <Switch>
          <Route exact path="/login">
            <Login />
          </Route>
          <Route exact path="/register">
            <Register />
          </Route>
        </Switch>
      </div>
      <div className={styles.picasso}>
        <h2 className={styles.tagline}>
          <BellIcon className={styles.bell} />
          Become
        </h2>
        <h2 className={styles.tagline}>the all Knowing.</h2>
      </div>
    </div>
  );
};

export default LogReg;
