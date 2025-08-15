import axios from 'axios';

const apiClient = axios.create({
  baseURL: "http://api.saarathi.com:8080/api/v1",
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
})

apiClient.interceptors.request.use(
  (config) => {
    return config
  }
)

let isRefreshing = false;
let failedQueue: Array<{
  resolve: (value?: any) => void;
  reject: (err: any) => void
}> = [];

const processQueue = (error: any = null) => {
  failedQueue.forEach(p => error ? p.reject(error) : p.resolve());
  failedQueue = [];
}

apiClient.interceptors.response.use(
  (response) => response,

  async (error) => {
    const originalRequest = error.config;

    if (error.response.status == 401 && !originalRequest._retry) {
      if (isRefreshing) {
        // queue the request until refresh is done
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject })
        }).then(() => apiClient(originalRequest));
      }

      originalRequest._retry = true
      isRefreshing = true

      try {
        await apiClient.post("/tokens/refresh")
        processQueue()
        return apiClient(originalRequest);
      } catch (refreshError) {
        processQueue(refreshError)
        if (typeof window !== "undefined") {
          window.location.href = "/rider/sign-in";
        }

        return Promise.reject(refreshError)
      } finally {
        isRefreshing = false
      }
    }

    return Promise.reject(error)
  }
)

export default apiClient;
