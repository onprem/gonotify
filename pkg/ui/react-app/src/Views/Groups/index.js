import React from 'react';

import Button from 'Components/Button';
import { ReactComponent as SendIcon } from 'Assets/icons/send.svg';

import styles from './groups.module.css';

const Card = ({ name, id, groups }) => {
  return (
    <div className={styles.card}>
      <div className={styles.up}>
        <h3 className={styles.cardName}>{name}</h3>
        <p className={styles.cardDetail}>ID: {id}</p>
        <p className={styles.cardDetail}>Members: {groups}</p>
      </div>
      <Button className={styles.btn}>
        SEND <SendIcon />
      </Button>
    </div>
  );
};

const Groups = () => {
  const all = [
    {
      id: 1,
      name: 'default',
      numbers: [{ id: 1 }, { id: 2 }],
    },
  ];

  return (
    <div className={styles.groups}>
      <h2 className={styles.heading}>All Groups</h2>
      <div className={styles.cards}>
        {all.map((g) => (
          <Card key={g.id} name={g.name} id={g.id} groups={g.numbers.length} />
        ))}
      </div>
    </div>
  );
};

export default Groups;
