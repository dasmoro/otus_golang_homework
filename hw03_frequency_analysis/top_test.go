package hw03_frequency_analysis //nolint:golint

import (
	"testing"

	"io/ioutil"
	"log"
	"net/http"

	"github.com/stretchr/testify/require"
)

// Change to true if needed
var taskWithAsteriskIsCompleted = true

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

var enText = `This book is largely concerned with Hobbits, and from its pages a reader may discover much of their character and a little of their history. Further information will also be found in the selection from the Red Book of Westmarch that has already been published, under the title of The Hobbit. That story was derived from the earlier chapters of the Red Book, composed by Bilbo himself, the first Hobbit to become famous in the world at large, and called by him There and Back Again, since they told of his journey into the East and his return: an adventure which later involved all the Hobbits in the great events of that Age that are here related.

Many, however, may wish to know more about this remarkable people from the outset, while some may not possess the earlier book. For such readers a few notes on the more important points are here collected from Hobbit-lore, and the first adventure is briefly recalled.

Hobbits are an unobtrusive but very ancient people, more numerous formerly than they are today; for they love peace and quiet and good tilled earth: a well-ordered and well-farmed countryside was their favourite haunt. They do not and did not understand or like machines more complicated than a forge-bellows, a water-mill, or a hand-loom, though they were skilful with tools. Even in ancient days they were, as a rule, shy of 'the Big Folk', as they call us, and now they avoid us with dismay and are becoming hard to find. They are quick of hearing and sharp-eyed, and though they are inclined to be fat and do not hurry unnecessarily, they are nonetheless nimble and deft in their movements. They possessed from the first the art of disappearing swiftly and silently, when large folk whom they do not wish to meet come blundering by; and this an they have developed until to Men it may seem magical. But Hobbits have never, in fact, studied magic of any kind, and their elusiveness is due solely to a professional skill that heredity and practice, and a close friendship with the earth, have rendered inimitable by bigger and clumsier races.

For they are a little people, smaller than Dwarves: less tout and stocky, that is, even when they are not actually much shorter. Their height is variable, ranging between two and four feet of our measure. They seldom now reach three feet; but they hive dwindled, they say, and in ancient days they were taller. According to the Red Book, Bandobras Took (Bullroarer), son of Isengrim the Second, was four foot five and able to ride a horse. He was surpassed in all Hobbit records only by two famous characters of old; but that curious matter is dealt with in this book.

As for the Hobbits of the Shire, with whom these tales are concerned, in the days of their peace and prosperity they were a merry folk. They dressed in bright colours, being notably fond of yellow and green; but they seldom wore shoes, since their feet had tough leathery soles and were clad in a thick curling hair, much like the hair of their heads, which was commonly brown. Thus, the only craft little practised among them was shoe-making; but they had long and skilful fingers and could make many other useful and comely things. Their faces were as a rule good-natured rather than beautiful, broad, bright-eyed, red-cheeked, with mouths apt to laughter, and to eating and drinking. And laugh they did, and eat, and drink, often and heartily, being fond of simple jests at all times, and of six meals a day (when they could get them). They were hospitable and delighted in parties, and in presents, which they gave away freely and eagerly accepted.

It is plain indeed that in spite of later estrangement Hobbits are relatives of ours: far nearer to us than Elves, or even than Dwarves. Of old they spoke the languages of Men, after their own fashion, and liked and disliked much the same things as Men did. But what exactly our relationship is can no longer be discovered. The beginning of Hobbits lies far back in the Elder Days that are now lost and forgotten. Only the Elves still preserve any records of that vanished time, and their traditions are concerned almost entirely with their own history, in which Men appear seldom and Hobbits are not mentioned at all. Yet it is clear that Hobbits had, in fact, lived quietly in Middle-earth for many long years before other folk became even aware of them. And the world being after all full of strange creatures beyond count, these little people seemed of very little importance. But in the days of Bilbo, and of Frodo his heir, they suddenly became, by no wish of their own, both important and renowned, and troubled the counsels of the Wise and the Great.

Those days, the Third Age of Middle-earth, are now long past, and the shape of all lands has been changed; but the regions in which Hobbits then lived were doubtless the same as those in which they still linger: the North-West of the Old World, east of the Sea. Of their original home the Hobbits in Bilbo's time preserved no knowledge. A love of learning (other than genealogical lore) was far from general among them, but there remained still a few in the older families who studied their own books, and even gathered reports of old times and distant lands from Elves, Dwarves, and Men. Their own records began only after the settlement of the Shire, and their most ancient legends hardly looked further back than their Wandering Days. It is clear, nonetheless, from these legends, and from the evidence of their peculiar words and customs, that like many other folk Hobbits had in the distant past moved westward. Their earliest tales seem to glimpse a time when they dwelt in the upper vales of Anduin, between the eaves of Greenwood the Great and the Misty Mountains. Why they later undertook the hard and perilous crossing of the mountains into Eriador is no longer certain. Their own accounts speak of the multiplying of Men in the land, and of a shadow that fell on the forest, so that it became darkened and its new name was Mirkwood.

Before the crossing of the mountains the Hobbits had already become divided into three somewhat different breeds: Harfoots, Stoors, and Fallohides. The Harfoots were browner of skin, smaller, and shorter, and they were beardless and bootless; their hands and feet were neat and nimble; and they preferred highlands and hillsides. The Stoors were broader, heavier in build; their feet and hands were larger, and they preferred flat lands and riversides. The Fallohides were fairer of skin and also of hair, and they were taller and slimmer than the others; they were lovers of trees and of woodlands.`

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
	})

	t.Run("positive test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{"он", "а", "и", "что", "ты", "не", "если", "то", "его", "кристофер", "робин", "в"}
			require.Subset(t, expected, Top10(text))
		} else {
			expected := []string{"он", "и", "а", "что", "ты", "не", "если", "-", "то", "Кристофер"}
			require.ElementsMatch(t, expected, Top10(text))
		}
	})

	t.Run("English test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{"and", "the", "of", "they", "in", "their", "a", "were", "are", "that", "hobbits", "to"}
			require.Subset(t, expected, Top10(enText))
		} else {
			expected := []string{"and", "the", "of", "they", "in", "their", "a", "were", "are", "that", "hobbits", "to"}
			require.ElementsMatch(t, expected, Top10(enText))
		}
	})

	t.Run("Crazy", func(t *testing.T) {
		resp, err := http.Get("https://code.jquery.com/jquery-3.5.1.js")
		jquery, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		expected := []string{"if", "elem", "function", "return", "jquery", "the", "this", "i", "a", "0", "type", "var"}
		require.Subset(t, expected, Top10(string(jquery)))
	})
}
