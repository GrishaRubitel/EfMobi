-- CREATE TABLE public.artist_data (
-- 	artist varchar NOT NULL,
-- 	spotify_id varchar NULL,
-- 	spotify_link varchar NULL,
-- 	youtube_link varchar NULL,
-- 	genres jsonb NULL,
-- 	CONSTRAINT artist_data_pk PRIMARY KEY (artist)
-- );

-- CREATE TABLE public.sound_data (
-- 	id serial4 NOT NULL,
-- 	artist varchar NULL,
-- 	title varchar NULL,
-- 	release_date date NULL,
-- 	lyrics text NULL,
-- 	spotify_link varchar NULL,
-- 	spotify_id varchar NULL,
-- 	video_link varchar NULL,
-- 	CONSTRAINT sound_data_pk PRIMARY KEY (id),
-- 	CONSTRAINT sound_data_unique_sp_id UNIQUE (spotify_id),
-- 	CONSTRAINT sound_data_unique_sp_link UNIQUE (spotify_link),
-- 	CONSTRAINT sound_data_unique_vid_link UNIQUE (video_link),
-- 	CONSTRAINT sound_data_artist_data_fk FOREIGN KEY (artist) REFERENCES public.artist_data(artist) ON DELETE CASCADE ON UPDATE CASCADE
-- );

-- CREATE TABLE public.lib_user (
-- 	id serial NOT NULL,
-- 	nickname varchar NOT NULL,
-- 	spotify_link varchar NULL,
-- 	CONSTRAINT lib_user_pk PRIMARY KEY (id),
-- 	CONSTRAINT lib_user_unique UNIQUE (nickname),
-- 	CONSTRAINT lib_user_unique_1 UNIQUE (spotify_link)
-- );

-- CREATE TYPE access_level AS ENUM ('private', 'public_read', 'public_edit');

-- CREATE TABLE public.playlist_data (
-- 	list_id int NOT NULL,
-- 	user_id int NULL,
-- 	list_name varchar NULL,
-- 	access_level access_level NULL,
-- 	CONSTRAINT playlist_data_pk PRIMARY KEY (list_id),
-- 	CONSTRAINT playlist_data_lib_user_fk FOREIGN KEY (user_id) REFERENCES public.lib_user(id) ON DELETE CASCADE ON UPDATE CASCADE
-- );

-- CREATE TABLE public.playlist_tracks (
-- 	list_id int4 NULL,
-- 	track_id int4 NULL,
-- 	CONSTRAINT playlist_tracks_pk PRIMARY KEY (list_id,track_id),
-- 	CONSTRAINT playlist_tracks_playlist_data_fk FOREIGN KEY (list_id) REFERENCES public.playlist_data(list_id) ON DELETE CASCADE ON UPDATE CASCADE,
-- 	CONSTRAINT playlist_tracks_sound_data_fk FOREIGN KEY (track_id) REFERENCES public.sound_data(id) ON DELETE CASCADE ON UPDATE CASCADE
-- );

-- CREATE INDEX idx_playlist_track ON public.playlist_tracks (list_id, track_id);

-- CREATE INDEX idx_sound_title ON public.sound_data (artist, title);


-- INSERT INTO public.lib_user (nickname, spotify_link)
-- VALUES 
--     ('GeshkaGorin', 'https://spotify.com/user/user1'),
--     ('OlegGazmanov', 'https://spotify.com/user/user2'),
--     ('GrishaRubitel', 'https://open.spotify.com/user/31rj2ijltqclk33ohfrvvf4llpbm');
   
-- INSERT INTO public.playlist_data (list_id, user_id, list_name, access_level)
-- VALUES
--     (1, 1, 'Geshkina muzyaks', 'private'),
--     (2, 1, 'Russian rock', 'public_read'),
--     (3, 2, 'Officers', 'private'),
--     (4, 2, 'Esual', 'public_edit'),
--     (5, 3, 'Favorite music', 'public_read');

-- INSERT INTO public.playlist_tracks (list_id, track_id)
-- VALUES
--     (1, 1),  
--     (1, 2),  
--     (1, 3),  
--     (2, 4),  
--     (2, 5),  
--     (3, 6),  
--     (3, 7),  
--     (4, 8),  
--     (5, 9),  
--     (5, 10); 

--------------------
--- NEW DATABASE ---
--------------------


CREATE TABLE public.artist_data (
	artist varchar NOT NULL,
	spotify_id varchar NULL,
	spotify_link varchar NULL,
	youtube_link varchar NULL,
	genres jsonb NULL,
	CONSTRAINT artist_data_pk PRIMARY KEY (artist)
);

CREATE TABLE public.sound_data (
	id serial4 NOT NULL,
	artist varchar NULL,
	title varchar NULL,
	release_date date NULL,
	lyrics text NULL,
	spotify_link varchar NULL,
	spotify_id varchar NULL,
	video_link varchar NULL,
	CONSTRAINT sound_data_pk PRIMARY KEY (id),
	CONSTRAINT sound_data_unique_sp_id UNIQUE (spotify_id),
	CONSTRAINT sound_data_unique_sp_link UNIQUE (spotify_link),
	CONSTRAINT sound_data_unique_vid_link UNIQUE (video_link),
	CONSTRAINT sound_data_artist_data_fk FOREIGN KEY (artist) REFERENCES public.artist_data(artist) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE INDEX idx_sound_title ON public.sound_data USING btree (artist, title);

INSERT INTO public.artist_data (artist, spotify_id, spotify_link, youtube_link, genres)
VALUES 
    ('JPEGMafia', '6yJ6QQ3Y5l0s0tn7b0arrO', 'https://open.spotify.com/artist/6yJ6QQ3Y5l0s0tn7b0arrO?si=ZGaEzMFvSgODUF_f7lJc3A', 'https://www.youtube.com/@JPEGMAFIA', '["experimental hip hop"]'::jsonb),
    ('Slipknot', '05fG473iIaoy82BF1aGhL8', 'https://open.spotify.com/artist/05fG473iIaoy82BF1aGhL8', 'https://youtube.com/@slipknot', '["metal", "nu metal"]'::jsonb),
    ('Radiohead', '4Z8W4fKeB5YxbusRsdQVPb', 'https://open.spotify.com/artist/4Z8W4fKeB5YxbusRsdQVPb', 'https://youtube.com/@radiohead', '["alternative rock", "art rock"]'::jsonb),
    ('Twenty One Pilots', '3YQKmKGau1PzlVlkL1iodx', 'https://open.spotify.com/artist/3YQKmKGau1PzlVlkL1iodx', 'https://youtube.com/@twentyonepilots', '["alternative", "pop rock"]'::jsonb),
    ('TV Girl', '0Y6dVaC9DZtPNH4591M42W', 'https://open.spotify.com/artist/0Y6dVaC9DZtPNH4591M42W', 'https://youtube.com/@teeveegirl', '["indie pop"]'::jsonb),
    ('Linkin Park', '6XyY86QOPPrYVGvF9ch6wz', 'https://open.spotify.com/artist/6XyY86QOPPrYVGvF9ch6wz?si=lgh0bUyeS9upOivlAkNzVA', 'https://youtube.com/@LinkinPark', '["nu metal", "alternative rock"]'::jsonb),
    ('Korn', '3RNrq3jvMZxD9ZyoOZbQOD', 'https://open.spotify.com/artist/3RNrq3jvMZxD9ZyoOZbQOD', 'https://youtube.com/@kornchannel', '["metal", "nu metal"]'::jsonb),
    ('Denzel Curry', '6fxyWrfmjcbj5d12gXeiNV', 'https://open.spotify.com/artist/6fxyWrfmjcbj5d12gXeiNV', 'https://youtube.com/@DENZELCURRYPH', '["hip hop"]'::jsonb),
    ('Kali Uchis', '1U1el3k54VvEUzo3ybLPlM', 'https://open.spotify.com/artist/1U1el3k54VvEUzo3ybLPlM', 'https://youtube.com/@KALIUCHIS', '["alternative R&B", "neo soul"]'::jsonb),
    ('Gorillaz', '3AA28KZvwAUcZuOKwyblJQ', 'https://open.spotify.com/artist/3AA28KZvwAUcZuOKwyblJQ', 'https://youtube.com/@Gorillaz', '["alternative", "virtual band"]'::jsonb);

INSERT INTO public.sound_data (id, artist, title, release_date, lyrics, spotify_link, spotify_id, video_link) VALUES(1, 'JPEGMafia', '1539 N.Calvert', '2021-09-08', '...', 'https://open.spotify.com/track/6XyxCBp6x3jvtxXvMN5sAA', '6XyxCBp6x3jvtxXvMN5sAA', 'https://www.youtube.com/watch?v=PO3mri47s7M');
INSERT INTO public.sound_data (id, artist, title, release_date, lyrics, spotify_link, spotify_id, video_link) VALUES(2, 'Slipknot', 'Dead Memories', '2008-09-30', '...', 'https://open.spotify.com/track/0HAr4QR1xI8nwC7VfzYidu', '0HAr4QR1xI8nwC7VfzYidu', 'https://www.youtube.com/watch?v=9gsAz6S_zSw	');
INSERT INTO public.sound_data (id, artist, title, release_date, lyrics, spotify_link, spotify_id, video_link) VALUES(3, 'Radiohead', 'All I Need', '2007-10-10', '...', 'https://open.spotify.com/track/5Qv2Nby1xTr9pQyjkrc94J', '5Qv2Nby1xTr9pQyjkrc94J', 'https://www.youtube.com/watch?v=wUL8NklXDsw');
INSERT INTO public.sound_data (id, artist, title, release_date, lyrics, spotify_link, spotify_id, video_link) VALUES(4, 'Twenty One Pilots', 'Heavy Dirty Soul', '2015-05-19', '...', 'https://open.spotify.com/track/7i9763l5SSfOnqZ35VOcfy?si=97a4d7022a97410c', '7i9763l5SSfOnqZ35VOcfy', 'https://www.youtube.com/watch?v=r_9Kf0D5BTs');
INSERT INTO public.sound_data (id, artist, title, release_date, lyrics, spotify_link, spotify_id, video_link) VALUES(5, 'TV Girl', 'Lovers Rock', '2014-06-05', '...', 'https://open.spotify.com/track/6dBUzqjtbnIa1TwYbyw5CM', '6dBUzqjtbnIa1TwYbyw5CM', 'https://www.youtube.com/watch?v=j_sG_Juncn8');
INSERT INTO public.sound_data (id, artist, title, release_date, lyrics, spotify_link, spotify_id, video_link) VALUES(6, 'Linkin Park', 'Lying From You', '2003-03-25', '...', 'https://open.spotify.com/track/4qVR3CF8FuFvHN4L6vXlB1', '4qVR3CF8FuFvHN4L6vXlB1', 'https://www.youtube.com/watch?v=NjdgcHdzvac');
INSERT INTO public.sound_data (id, artist, title, release_date, lyrics, spotify_link, spotify_id, video_link) VALUES(7, 'Korn', 'Good God', '1996-08-18', '...', 'https://open.spotify.com/track/5JrajjztyjvkuUB8ZqzUML', '5JrajjztyjvkuUB8ZqzUML', 'https://www.youtube.com/watch?v=GHkUCSeTi2I');
INSERT INTO public.sound_data (id, artist, title, release_date, lyrics, spotify_link, spotify_id, video_link) VALUES(8, 'Denzel Curry', 'Walkin', '2022-01-24', '...', 'https://open.spotify.com/track/1q8DwZtQen5fvyB7cKbShC', '1q8DwZtQen5fvyB7cKbShC', 'https://www.youtube.com/watch?v=fOO1mWLGhh8');
INSERT INTO public.sound_data (id, artist, title, release_date, lyrics, spotify_link, spotify_id, video_link) VALUES(9, 'Kali Uchis', 'Igual Que Un Angel', '2020-12-04', '...', 'https://open.spotify.com/track/6XaJfhwof7qIgbbXO5tIQI', '6XaJfhwof7qIgbbXO5tIQI', 'https://www.youtube.com/watch?v=YR1t_MUN8I4');
INSERT INTO public.sound_data (id, artist, title, release_date, lyrics, spotify_link, spotify_id, video_link) VALUES(10, 'Gorillaz', 'She Is My Collar', '2017-04-28', '...', 'https://open.spotify.com/track/3lIxtCaROdRDuTnNBDm3n2', '3lIxtCaROdRDuTnNBDm3n2', 'https://www.youtube.com/watch?v=vq_5126alC8');
INSERT INTO public.sound_data (id, artist, title, release_date, lyrics, spotify_link, spotify_id, video_link) VALUES(11, 'Kali Uchis', 'Dead To Me', '2018-04-06', 'You''re dead to me
You''re obsessed, just let me go
You''re dead to me
I''m not somebody you know
You''re dead to me
Could you just leave me alone?
You''re dead to me', 'https://open.spotify.com/track/6LOZws7T3jqZz78unPgFF9', '6LOZws7T3jqZz78unPgFF9', 'https://www.youtube.com/watch?v=OcUDK4kAUIw');
