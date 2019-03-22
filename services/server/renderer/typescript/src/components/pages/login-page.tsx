import * as React from "react";
import {registerComponent} from "preact-rpc";
import {IPageProps, Page} from "../page";
import {InputGroup} from "../bootstrap/input-group";
import {AlertMessage} from "../bootstrap/alert-message";

interface IProps extends IPageProps {

}

export class LoginPage extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <Page authPayload={this.props.authPayload} id={"login-page"}>
            <div className={"row"}>
                <div className={"col-sm-4 offset-sm-4"}>
                    <h1 className={"text-center"}>login</h1>
                    <AlertMessage type={"danger"} message={this.props.error} />
                    <form method={"POST"} action={"/user/login"}>
                        <InputGroup name={"username"} type={"text"} placeholder={"username"}/>
                        <InputGroup name={"password"} type={"password"} placeholder={"password"}/>
                        <button className="btn btn-primary btn-block" type="submit">login</button>
                    </form>
                </div>
            </div>
        </Page>;
    }

}

registerComponent('LoginPage', LoginPage);
