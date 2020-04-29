export default ( headers = {}, endpoint, method = 'GET', body = '') => {
  console.log('fetching', endpoint, body)
  let options = {
    method,
    headers,
  };

  if (method === 'POST') {
    options = {
      ...options,
      headers: {
        ...options?.headers,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body),
    };
  }

  return fetch(endpoint, options).then((res) => res.json());
};
