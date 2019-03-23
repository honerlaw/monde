import * as React from "react";
import {IPageNavbarProps, PageNavbar} from "./page-navbar";
import {AlertMessage} from "./bootstrap/alert-message";

export interface IPageProps extends IPageNavbarProps {
    id: string;
    error?: string;
}

export class Page extends React.Component<IPageProps, {}> {

    public render(): JSX.Element {
        return <html>
        <head>
            <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
            <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1"/>
            <title>vueon</title>
            <link rel={"stylesheet"} href={"/css/index.css"}/>
        </head>
        <body>
        <PageNavbar authPayload={this.props.authPayload}/>
        <div className="container" id={this.props.id}>
            <AlertMessage type={"danger"} message={this.props.error} />
            {this.props.children}
        </div>
        </body>
        </html>;
    }

}