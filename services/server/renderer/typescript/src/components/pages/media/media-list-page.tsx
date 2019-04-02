import * as React from "react";
import "./media-list-page.scss";
import {IGlobalProps} from "../../../global-props";
import {IUploadForm} from "../../media/upload-form";
import {MediaPendingItem} from "./media-list-page/media-pending-item";
import {MediaItem} from "./media-list-page/media-item";

export interface IMediaVideoResponse {
    type: string;
    width: number;
    height: number;
    url: string;
}

export interface IMediaResponse {
    id: string;
    title: string;
    description: string;
    transcoding_status: string;
    hashtags: string[];
    is_published: boolean;
    can_publish: boolean;
    thumbnails: string[];
    videos: IMediaVideoResponse[];
}

interface IProps extends IGlobalProps {
    uploads: IMediaResponse[];
    uploadForm: IUploadForm;
}

/**
 * @todo
 * - display thumbnail after upload so the user can see what it is
 */
export class MediaListPage extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <div id={"upload-list-page"}>
            <ol className={"upload-list"}>
                {this.props.uploads.map((upload: IMediaResponse): JSX.Element => {
                    if (upload.transcoding_status !== "Complete") {
                        return <MediaPendingItem key={upload.id} status={upload.transcoding_status}/>;
                    }
                    return <MediaItem key={upload.id} upload={upload}/>;
                })}
            </ol>
        </div>;
    }

}
