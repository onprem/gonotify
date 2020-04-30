import React from 'react';
import useSWR from 'swr';

import styles from './notif.module.css';

const NotifList = ({ notifs }) => {
  const list = notifs.map((n) => {
    const t = new Date(n.timeSt);
    return (
    <div key={n.id} className={styles.row}>
      <span className={styles.idSpan}>{n.id}</span>
      <span className={styles.groupSpan}>{n.groupName}</span>
      <span className={styles.timeSpan}>{t.toLocaleString()}</span>
      <span className={styles.bodySpan}>{n.body}</span>
    </div>
  )});
  return (
    <div className={styles.list}>
      <div className={styles.topRow}>
        <span className={styles.idSpan}>ID</span>
        <span className={styles.groupSpan}>Group</span>
        <span className={styles.timeSpan}>TimeStamp</span>
        <span className={styles.bodySpan}>Body</span>
      </div>
      {list}
    </div>
  );
};

const Notifications = ({ groups }) => {
  const { data, error } = useSWR('/api/v1/notifications');

  if (error) return <h2>Some error occured.</h2>;
  if (!data) return <h2>Loading...</h2>;

  return (
    <div className={styles.notif}>
      <h2 className={styles.heading}>All Notifications</h2>
      <NotifList notifs={data.notifications} />
    </div>
  );
};

export default Notifications;
