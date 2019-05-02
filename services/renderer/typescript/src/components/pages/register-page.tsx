import * as React from "react";
import {InputGroup} from "../bootstrap/input-group";
import {AlertMessage} from "../bootstrap/alert-message";
import {IGlobalProps} from "../../global-props";
import "./register-page.scss";

interface IProps extends IGlobalProps {
    email?: string;
}

export class RegisterPage extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <div id={"register-page"}>
            <div className={"row"}>
                <div className={"col-sm-4 offset-sm-4"}>
                    <h1 className={"text-center"}>register</h1>
                    <AlertMessage type={"danger"} message={this.props.error} />
                    <form method={"POST"} action={"/user/register"}>
                        <InputGroup name={"email"} type={"text"} value={this.props.email} placeholder={"email"}/>
                        <InputGroup name={"password"} type={"password"} placeholder={"password"}/>
                        <InputGroup name={"verify_password"} type={"password"} placeholder={"verify password"}/>
                        <button className="btn btn-primary btn-block" type="submit">register</button>
                    </form>
                </div>
            </div>
        </div>;
    }

}
