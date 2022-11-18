
import { Link ,useParams} from "react-router-dom";
import axios from '../api/axios';


const VerifyScreen = () => {
    const { verificationCode } = useParams();
    axios.post('http://localhost:8080/v1/user/verifyemail', {
    "verificationCode": verificationCode,
    "userid": "string"
  });

    return (
        <section>
            <h1>Verify Page</h1>
            <h1>{verificationCode}</h1>

            <br />
            <br />
            <div className="flexGrow">
                <Link to="/">Home</Link>
            </div>
        </section>
    )
}

export default VerifyScreen
