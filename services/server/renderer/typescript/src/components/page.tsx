import * as React from "react";
import {IPageNavbarProps, PageNavbar} from "./page-navbar";

interface IProps extends IPageNavbarProps {
    options: string | null;
}

export class Page extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <html>
        <head>
            <meta httpEquiv={"Content-Type"} content={"text/html; charset=UTF-8"}/>
            <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1"/>
            <title>vueon</title>
            <link rel={"stylesheet"} href={"/css/index.css"}/>
        </head>
        <body>
        <PageNavbar authPayload={this.props.authPayload} uploadForm={this.props.uploadForm}/>
        <div className="container">
            {this.props.children}
        </div>

        <script type={"text/javascript"} suppressHydrationWarning={true}
                dangerouslySetInnerHTML={{ __html: `window.hydrateOptions = ${this.props.options}`}} />
        <script type={"text/javascript"} src={"/js/bundle.js"}/>
        </body>
        </html>;
    }

}