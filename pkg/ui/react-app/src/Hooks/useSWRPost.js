import { useState } from 'react';
import useSWR from 'swr';

const useSWRPost = (endpoint, swrOpts) => {
  const [values, runFetch] = useState();

  const swrOut = useSWR(values ? [endpoint, 'POST', values] : null, {
    revalidateOnFocus: false,
    ...swrOpts,
  });

  return [runFetch, swrOut];
};

export default useSWRPost;
