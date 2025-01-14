import axios from 'axios'

const NonTokenHTTP = axios.create({
  headers: {
    'Content-Type': 'application/json',
  },
})

NonTokenHTTP.interceptors.response.use(
  async (response) => {
    await sessionStorage.setItem('user', JSON.stringify(response.data))
    return response
  },
  (err) => {
    console.error(err)
  },
)

const TokenHTTP = axios.create()

TokenHTTP.interceptors.request.use(function (config) {
  config!.headers!.Authorization = `Bearer ${JSON.parse(sessionStorage.getItem('user')!)?.access}`
  return config
})

TokenHTTP.interceptors.response.use(
  (response) => response,
  async (err) => {
    const origConfig = err.config
    if (err?.response?.status === 401) {
      try {
        const username = sessionStorage.getItem('username')
        const password = sessionStorage.getItem('password')
        const dataBody = { username: username, password: password }
        const resp = await NonTokenHTTP.post('/server/token/pair', dataBody)
        sessionStorage.setItem('user', JSON.stringify(resp.data))
        return TokenHTTP(origConfig)
      } catch (error) {
        return Promise.reject(err)
      }
    }
  },
)
export { NonTokenHTTP, TokenHTTP }
