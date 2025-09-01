package db

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/cprakhar/gopher-social/internal/store"
)

func Seed(store store.Store) {
	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Printf("failed to create user %s: %v", user.Username, err)
			return
		}
	}

	posts := generatePosts(users, 300)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Printf("failed to create post %s: %v", post.Title, err)
			return
		}
	}

	comments := generateComments(users, posts, 500)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Printf("failed to create comment on post %s: %v", comment.PostID, err)
			return
		}
	}

	log.Println("Database seeding completed successfully")
}

var usernames = []string{
	"alice", "bob", "charlie", "dave", "eve",
	"frank", "grace", "heidi", "ivan", "judy",
	"mallory", "niaj", "oscar", "peggy", "quinn",
	"rudy", "sybil", "trent", "ursula", "victor",
	"wendy", "xander", "yvonne", "zach",
}

var titles = []string{
	"The Silent Voyager", "Echoes of Tomorrow", "Master of Shadows",
	"The Quantum Blueprint", "Whispers of the Cosmos", "The Last Cipher",
	"Architect of Dreams", "The Fractal Path", "Stormbound Horizons",
	"Keeper of Lost Code", "The Shattered Nexus", "Chronicles of the Void",
	"The Ember Scrolls", "Warden of Infinity", "The Forgotten Pulse",
	"Prisms of Eternity", "The Iron Nomad", "Songs of the Binary Dawn",
	"The Crystal Dominion", "Echoes in the Circuit",
}

var content = []string{
	"Exploring the hidden patterns in everyday data reveals how small decisions shape larger outcomes. Data often tells stories we don’t notice at first glance, and by connecting these threads, we discover meaningful insights. It’s in these patterns that innovation quietly begins.",
	"A journey through the cosmos of imagination allows us to explore infinite possibilities. The human mind is capable of constructing entire universes with just a thought. Imagination is not just an escape, but a tool that fuels discovery and progress.",
	"Unveiling the secrets behind great innovations is like peeling back the layers of history. Every invention is built on countless unseen trials, failures, and breakthroughs. What looks like overnight success is often the product of years of persistence.",
	"When technology meets creativity, magic happens in ways that transform the world. Artists and engineers together create solutions that are both functional and beautiful. The intersection of art and science is where progress thrives.",
	"The art of simplifying complex problems lies in breaking them into smaller, manageable steps. Great problem-solvers don’t chase complexity; they embrace clarity. Simplicity often turns out to be the most elegant solution.",
	"A story of resilience and determination often begins with failure. Every obstacle faced becomes a stepping stone toward growth. True success belongs to those who keep moving forward despite the setbacks.",
	"Capturing moments that define our lives allows us to revisit emotions long after the moment has passed. Whether through photography, writing, or memory, these fragments of time remind us of who we are. They shape our identities and connect us to others.",
	"The future belongs to those who dream big and refuse to settle. Visionaries imagine worlds not yet built and chase them relentlessly. Every leap forward in history started with a dream others thought impossible.",
	"Breaking boundaries with every new discovery is what propels humanity forward. Science, art, and technology all expand when someone dares to ask, 'What if?'. Boundaries only exist until someone has the courage to cross them.",
	"A spark of curiosity can ignite revolutions in thought and innovation. Many of the greatest breakthroughs began with a simple question. Curiosity keeps us moving, learning, and creating beyond what we already know.",
	"Learning from failures, building for success, is a cycle that repeats endlessly in progress. Failure is not the end but a necessary teacher. Those who embrace mistakes as lessons unlock the path to lasting success.",
	"The rhythm of progress never stops, even when we don’t notice it. Small, incremental changes often build into massive transformations over time. Progress hums quietly in the background until it suddenly becomes visible.",
	"Ideas are seeds that grow into change, but only when nurtured with effort. A single idea can alter entire industries, cultures, or even civilizations. The real magic lies in execution, not just inspiration.",
	"Innovation thrives where passion lives, because passion sustains the energy needed to keep going. Without genuine care, progress falters. Passionate people create movements that shift how the world operates.",
	"Exploring the nexus between man and machine opens doors to new possibilities. Artificial intelligence, robotics, and automation redefine what it means to create and collaborate. Together, they extend the reach of human potential.",
	"Every step forward shapes the future, no matter how small. Incremental improvements accumulate into revolutions that change history. Progress is rarely dramatic—it is built step by step, day by day.",
	"Crafting solutions for a better tomorrow requires courage and creativity. The challenges we face today are opportunities to build something extraordinary. Solutions designed with care today will shape the lives of generations tomorrow.",
	"Inspiration comes from the simplest things, often hidden in plain sight. A word, a moment, or a fleeting encounter can spark a cascade of new thoughts. The beauty of inspiration is that it arrives when we least expect it.",
	"The digital frontier is just the beginning of what humanity can achieve. Every advancement in technology opens up more doors, revealing possibilities that once seemed unreachable. This frontier is not an end but a continuous expansion.",
	"Stories connect us across time and space, reminding us of our shared humanity. Through storytelling, knowledge and emotion pass from one generation to another. They preserve who we are while guiding us toward who we might become.",
}

var tags = []string{
	"technology", "innovation", "science", "art", "creativity",
	"future", "progress", "inspiration", "curiosity", "resilience",
	"discovery", "imagination", "design", "development", "research",
	"data", "learning", "growth", "change", "vision",
}

var comments = []string{
	"Great post! Really made me think.",
	"I completely agree with your points.",
	"This is a fascinating perspective.",
	"Thanks for sharing your insights!",
	"I learned something new today.",
	"Can't wait to read more from you.",
	"Your writing is very engaging.",
	"This topic is so relevant right now.",
	"Well said! I appreciate your thoughts.",
	"Looking forward to your next post.",
	"This was very informative, thank you!",
	"I love how you explained this concept.",
	"Such a unique take on the subject!",
	"You've inspired me to explore this further.",
	"Fantastic read, I enjoyed every bit of it.",
	"This really resonated with me.",
	"Your examples were very helpful.",
	"I appreciate the depth of your analysis.",
	"This has given me a lot to consider.",
	"Keep up the great work!",
}

func generateUsers(n int) []*store.User {
	users := make([]*store.User, n)
	for i := range n {
		username := usernames[i%len(usernames)] + fmt.Sprintf("%d", i/len(usernames)+1)
		email := username + "@example.com"
		users[i] = &store.User{
			Username: username,
			Email:    email,
			Password: []byte("password"), // In real scenarios, ensure to hash passwords
		}
	}
	return users
}

func generatePosts(users []*store.User, n int) []*store.Post {
	posts := make([]*store.Post, n)
	for i := range n {
		author := users[rand.IntN(len(users))]
		posts[i] = &store.Post{
			Title:    titles[rand.IntN(len(titles))],
			Content:  content[rand.IntN(len(content))],
			AuthorID: author.ID,
			Tags: []string{
				tags[rand.IntN(len(tags))],
				tags[rand.IntN(len(tags))],
			},
		}
	}
	return posts
}

func generateComments(users []*store.User, posts []*store.Post, n int) []*store.Comment {
	commentsList := make([]*store.Comment, n)
	for i := range n {
		author := users[rand.IntN(len(users))]
		post := posts[rand.IntN(len(posts))]
		commentsList[i] = &store.Comment{
			PostID:   post.ID,
			AuthorID: author.ID,
			Content:  comments[rand.IntN(len(comments))],
		}
	}
	return commentsList
}
