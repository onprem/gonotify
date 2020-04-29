import toast from 'cogo-toast';

const defOpts = {
  position: 'top-right',
}

export const error = (msg, options) => {
  toast.error(msg, {
    ...defOpts,
    ...options,
  });
}

export const success = (msg, options) => {
  toast.success(msg, {
    ...defOpts,
    ...options,
  });
}

export const info = (msg, options) => {
  toast.info(msg, {
    ...defOpts,
    ...options,
  });
}

export const loading = (msg, options) => {
  toast.loading(msg, {
    ...defOpts,
    ...options,
  });
}

export const warn = (msg, options) => {
  toast.warn(msg, {
    ...defOpts,
    ...options,
  });
}

export default {
  success,
  error,
  warn,
  loading,
  info,
}
