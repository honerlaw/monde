import * as React from "react";
import {ICommentResponse} from "../media-view-page";
import {CommentItem} from "./comment-item";

interface IProps {
    comments: ICommentResponse[];
}

export class CommentList extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <div id={"comment-container"}>
            {this.props.comments.map((comment: ICommentResponse): JSX.Element => {
                return <CommentItem comment={comment}/>
            })}
        </div>
    }

}