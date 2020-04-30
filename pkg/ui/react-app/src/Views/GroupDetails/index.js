import React from 'react';
import { useParams } from 'react-router-dom';

import Button from 'Components/Button';
import { ReactComponent as TrashIcon } from 'Assets/icons/trash.svg';

import styles from './groupdetails.module.css';

const Card = ({ node }) => {
  const { id, phone, numberID, lastMsgReceived } = node;

  return (
    <div className={styles.card}>
      <div className={styles.content}>
        <h3 className={styles.cardName}>{phone}</h3>
        <p className={styles.cardDetail}>ID: {id}</p>
        <p className={styles.cardDetail}>NumberID: {numberID}</p>
      </div>
      <Button className={styles.btn}>
        Remove <TrashIcon />
      </Button>
    </div>
  );
};

const GroupDetails = ({ groups }) => {
  const { name } = useParams();

  const [group] = groups.filter((g) => g.name === name.toLowerCase());

  if (!group) return <h2>Invalid Group name - {name}</h2>;

  return (
    <div className={styles.main}>
      <h2 className={styles.heading}>{group.name}</h2>
      <h3 className={styles.sub}>ID:&nbsp; {group.id}</h3>
      <div className={styles.cards}>
        {group.whatsappNodes.map((n) => (
          <Card key={n.id} node={n} />
        ))}
      </div>
    </div>
  );
};

export default GroupDetails;
