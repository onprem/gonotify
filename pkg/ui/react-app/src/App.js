import React, { useState, useEffect } from 'react';
import { Switch, Route, Redirect } from 'react-router-dom';
import { SWRConfig } from 'swr';

import fetcher from 'Utils/fetcher';
import { AuthProvider } from 'Context/auth';
import ProtectedRoute from 'Components/ProtectedRoute';

import Nav from 'Views/Nav';
import Home from 'Views/Home';
import LogReg from 'Views/LogReg';
import Dashboard from 'Views/Dashboard';

import styles from './App.module.css';

function App() {
  const [token, setToken] = useState('');

  useEffect(() => {
    if (!token) {
      const tk = localStorage.getItem('gonotify-token');
      if (tk) setToken(tk);
    }
  }, [token]);

  const logMeIn = (token) => {
    window.localStorage.setItem('gonotify-token', token);
    setToken(token);
  };

  const logMeOut = () => {
    window.localStorage.removeItem('gonotify-token');
    setToken('');
  };
  global.logOut = logMeOut;

  return (
    <AuthProvider value={{ token, setToken, logMeIn, logMeOut }}>
      <SWRConfig
        value={{
          fetcher: (...args) => fetcher({ Authorization: `Bearer ${token}` }, ...args),
          revalidateOnFocus: false,
        }}
      >
        <div className={styles.App}>
          <Route exact path={['/', '/login', '/register', '/verify/:phone']}>
            <Nav />
          </Route>
          <Switch>
            <Route exact path="/">
              <Home />
            </Route>
            <Route exact path={['/login', '/register', '/verify/:phone']}>
              <LogReg />
            </Route>
            <ProtectedRoute exact path="/dashboard">
              <Redirect to="/dashboard/groups" />
            </ProtectedRoute>
            <ProtectedRoute path="/dashboard">
              <Dashboard />
            </ProtectedRoute>
          </Switch>
        </div>
      </SWRConfig>
    </AuthProvider>
  );
}

export default App;
