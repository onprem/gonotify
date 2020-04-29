import React, { useState } from 'react';
import { Switch, Route } from 'react-router-dom';
import { SWRConfig } from 'swr';

import fetcher from 'Utils/fetcher';
import { AuthProvider } from 'Context/auth';

import Nav from 'Views/Nav';
import Home from 'Views/Home';
import LogReg from 'Views/LogReg';

import styles from './App.module.css';

function App() {
  const [token, setToken] = useState('');
  return (
    <AuthProvider value={{ token, setToken }}>
      <SWRConfig
        value={{
          fetcher: (...args) => fetcher({ Authorization: `Bearer ${token}` }, ...args),
        }}
      >
        <div className={styles.App}>
          <Nav />
          <Switch>
            <Route exact path="/">
              <Home />
            </Route>
            <Route exact path={['/login', '/register']}>
              <LogReg />
            </Route>
          </Switch>
        </div>
      </SWRConfig>
    </AuthProvider>
  );
}

export default App;
