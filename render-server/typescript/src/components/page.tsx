import * as React from "react";
import {IPageNavbarProps, PageNavbar} from "./page-navbar";

export interface IPageProps extends IPageNavbarProps {
    id: string;
    error?: string;
}

export class Page extends React.Component<IPageProps, {}> {

    public render(): JSX.Element {
        return <html>
        <head>
            <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1"/>
            <title>package</title>
            <link rel={"stylesheet"} href={"/css/index.css"}/>
        </head>
        <body>
        <PageNavbar authPayload={this.props.authPayload}/>
        <div className="container" id={this.props.id}>
            {this.props.children}
        </div>
        </body>
        </html>;
    }

}