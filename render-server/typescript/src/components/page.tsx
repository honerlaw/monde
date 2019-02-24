import * as React from "react";

export class Page extends React.Component<{}, {}> {

    public render(): JSX.Element {
        return <html>
        <head>
            <title>package</title>

            <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css"
                  integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm"
                  crossOrigin="anonymous" />
            <link rel={"stylesheet"} href={"/css/index.css"} />
        </head>
        <body>
            {this.props.children}
        </body>
        </html>;
    }

}