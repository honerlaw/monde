import * as React from "react";
import {InputGroup} from "../bootstrap/input-group";
import {AlertMessage} from "../bootstrap/alert-message";
import {IGlobalProps} from "../../global-props";
import "./login-page.scss";

interface IProps extends IGlobalProps {
    username?: string;
}

export class LoginPage extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <div id={"login-page"}>
            <div className={"row"}>
                <div className={"col-sm-4 offset-sm-4"}>
                    <h1 className={"text-center"}>login</h1>
                    <AlertMessage type={"danger"} message={this.props.error} />
                    <form method={"POST"} action={"/user/login"}>
                        <InputGroup name={"username"} type={"text"} value={this.props.username} placeholder={"username"}/>
                        <InputGroup name={"password"} type={"password"} placeholder={"password"}/>
                        <button className="btn btn-primary btn-block" type="submit">login</button>
                    </form>
                </div>
            </div>
        </div>;
    }

}

