import React from 'react';

import Nav from 'Views/Nav';
import Home from 'Views/Home';

import styles from './App.module.css';

function App() {
  return (
    <div className={styles.App}>
      <Nav />
      <Home />
    </div>
  );
}

export default App;
