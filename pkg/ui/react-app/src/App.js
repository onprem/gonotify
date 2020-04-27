import React from 'react';
import {ReactComponent as Logo} from './Assets/logo/logo.svg';
import styles from './App.module.css';

function App() {
  return (
    <div className={styles.App}>
      <Logo className={styles.logo} />
    </div>
  );
}

export default App;
