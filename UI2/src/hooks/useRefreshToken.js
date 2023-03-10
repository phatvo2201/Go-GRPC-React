import axios from 'axios';
import { useCookies } from 'react-cookie';



const useRefreshToken = () => {
    const [cookies, setCookie] = useCookies(['token','rftoken']);
    const refreshToken = cookies.rftoken

    const config = {
        headers: { Authorization: `Bearer ${refreshToken}`,'Content-Type': 'application/json' }
    }

    const refresh = async () => {
        //call to get new token by refresh token
        const response = await axios.post('http://localhost:8080/api/v1/auth/refresh',
            JSON.stringify({}),
            {
                headers: { 'Content-Type': 'application/json',Authorization: `Bearer ${refreshToken}` },
                withCredentials: true
            }
        );

    let accessToken = response.data.access_token
    setCookie('token',accessToken)


        return response.data.access_token;
    }
    return refresh;
};

export default useRefreshToken;
