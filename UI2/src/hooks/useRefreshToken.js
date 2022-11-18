import axios from 'axios';
import useAuth from './useAuth';
import { useCookies } from 'react-cookie';



const useRefreshToken = () => {
    const [cookies, setCookie] = useCookies(['token','rftoken']);

    
    const { auth ,setAuth } = useAuth();
    
    // const refreshToken = auth?.refreshToken
    const refreshToken = cookies.rftoken.refreshToken

    const config = {
        headers: { Authorization: `Bearer ${refreshToken}`,'Content-Type': 'application/json' }
    }

    const refresh = async () => {
        // const response = await axios.post("http://localhost:8080/v1/refresh",
        //     config
        // );
        // setAuth(prev => {
        //     console.log(JSON.stringify(prev));
        //     console.log(response.data.accessToken);
        //     return { ...prev, accessToken: response.data.accessToken }
        // });

        const response = await axios.post('http://localhost:8080/v1/refresh',
            JSON.stringify({}),
            {
                headers: { 'Content-Type': 'application/json',Authorization: `Bearer ${refreshToken}` },
                withCredentials: false
            }
        );

setAuth(prev => {
    console.log(JSON.stringify(prev));
    console.log("lolololol")
    console.log(response.data.access_token);
    return { ...prev, accessToken: response.data.access_token }
});

let accessToken = response.data.access_token
setCookie('token',{ accessToken })

//  setCookie('token',{ accessToken })
    
        return response.data.accessToken;
    }
    return refresh;
};

export default useRefreshToken;
