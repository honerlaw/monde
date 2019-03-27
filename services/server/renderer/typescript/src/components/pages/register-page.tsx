import * as React from "react";
import {registerComponent} from "preact-rpc";
import {IPageProps, Page} from "../page";
import {InputGroup} from "../bootstrap/input-group";
import {AlertMessage} from "../bootstrap/alert-message";

interface IProps extends IPageProps {
    username: string;
}

export class RegisterPage extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <Page authPayload={this.props.authPayload} id={"register-page"}>
            <div className={"row"}>
                <div className={"col-sm-4 offset-sm-4"}>
                    <h1 className={"text-center"}>register</h1>
                    <AlertMessage type={"danger"} message={this.props.error} />
                    <form method={"POST"} action={"/user/register"}>
                        <InputGroup name={"username"} type={"text"} value={this.props.username} placeholder={"username"}/>
                        <InputGroup name={"password"} type={"password"} placeholder={"password"}/>
                        <InputGroup name={"verify_password"} type={"password"} placeholder={"verify password"}/>
                        <button className="btn btn-primary btn-block" type="submit">register</button>
                    </form>
                </div>
            </div>
        </Page>;
    }

}

registerComponent('RegisterPage', RegisterPage);
