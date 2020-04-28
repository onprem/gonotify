import React from 'react';
import Button from 'Components/Button';

import {ReactComponent as HeroImg} from 'Assets/images/hero.svg';
import {ReactComponent as ArrowIcon} from 'Assets/icons/arrow.svg';
import styles from './home.module.css';

const Home = () => {
  return (
    <div className={styles.home}>
      <div className={styles.content}>
        <h1 className={styles.heading}>
          Staying up-to-date gets <span className={styles.accent}>more</span> easier
        </h1>
        <p className={styles.para}>
          Stay in control of all your products with seamless notifications on your favourite social
          platform.
        </p>
        <Button size="large" className={styles.btn}>JOIN NOW <ArrowIcon className={styles.arrow} /></Button>
      </div>
      <HeroImg className={styles.hero} />
    </div>
  );
};

export default Home;
